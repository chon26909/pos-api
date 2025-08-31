// app/service/order_service.go
package service

import (
	"errors"
	"pos-api/app/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func (s *OrderService) EnsureActiveSession(tableID uint64) (*model.TableSession, error) {
	var sess model.TableSession
	err := s.DB.Where("table_id = ? AND is_active = ?", tableID, true).First(&sess).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sess = model.TableSession{
			TableID:  tableID,
			Token:    uuid.New().String(),
			IsActive: true,
		}
		if err := s.DB.Create(&sess).Error; err != nil {
			return nil, err
		}
		// flip table status
		s.DB.Model(&model.Table{}).Where("id = ?", tableID).Update("status", "occupied")
		return &sess, nil
	}
	return &sess, err
}

func (s *OrderService) GetSessionByToken(token string) (*model.TableSession, error) {
	var sess model.TableSession
	if err := s.DB.Where("token = ? AND is_active = ?", token, true).Preload("Table").First(&sess).Error; err != nil {
		return nil, err
	}
	return &sess, nil
}

type OrderItemInput struct {
	ProductID uint64
	Qty       int
}

func (s *OrderService) CreateOrder(tableToken string, items []OrderItemInput, note *string) (*model.Order, error) {
	sess, err := s.GetSessionByToken(tableToken)
	if err != nil {
		return nil, err
	}

	// โหลดราคา products
	var products []model.Product
	var ids []uint64
	for _, it := range items {
		ids = append(ids, it.ProductID)
	}
	if err := s.DB.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	priceMap := make(map[uint64]float64)
	for _, p := range products {
		priceMap[p.ID] = p.Price
	}

	order := model.Order{
		TableSessionID: sess.ID,
		Status:         "preparing",
		Note:           note,
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var total float64 = 0
		for _, it := range items {
			unit := priceMap[it.ProductID]
			if unit == 0 {
				return errors.New("product not found or price zero")
			}
			oi := model.OrderItem{
				OrderID:   order.ID,
				ProductID: it.ProductID,
				Qty:       it.Qty,
				UnitPrice: unit,
				Status:    "preparing",
			}
			if err := tx.Create(&oi).Error; err != nil {
				return err
			}
			total += unit * float64(it.Qty)
		}
		return tx.Model(&order).Update("total_price", total).Error
	})

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) MarkCheckoutRequested(tableToken string) error {
	sess, err := s.GetSessionByToken(tableToken)
	if err != nil {
		return err
	}

	// แจ้งเตือนแอดมิน + set orders เป็น checkout_requested (เฉพาะที่ยังไม่ paid/cancelled)
	if err := s.DB.Create(&model.AdminAlert{
		TableSessionID: sess.ID,
		Type:           "checkout_request",
	}).Error; err != nil {
		return err
	}

	return s.DB.Model(&model.Order{}).
		Where("table_session_id = ? AND status NOT IN ('paid','cancelled')", sess.ID).
		Update("status", "checkout_requested").Error
}

func (s *OrderService) SetOrderStatus(orderID uint64, status string) error {
	return s.DB.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (s *OrderService) SettleAndClose(tableSessionID uint64) error {
	now := time.Now()
	// ปิด session และเคลียร์โต๊ะ
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).
			Where("table_session_id = ? AND status != 'paid'", tableSessionID).
			Update("status", "paid").Error; err != nil {
			return err
		}
		if err := tx.Model(&model.TableSession{}).
			Where("id = ?", tableSessionID).
			Updates(map[string]any{"is_active": false, "closed_at": &now}).Error; err != nil {
			return err
		}
		// flip table status -> available
		var ts model.TableSession
		if err := tx.First(&ts, tableSessionID).Error; err == nil {
			tx.Model(&model.Table{}).Where("id = ?", ts.TableID).Update("status", "available")
		}
		// mark alerts as read
		tx.Model(&model.AdminAlert{}).
			Where("table_session_id = ? AND is_read = false", tableSessionID).
			Update("is_read", true)
		return nil
	})
}
