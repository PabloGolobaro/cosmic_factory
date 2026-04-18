package service

import (
	"context"
	"slices"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// InventoryServer реализует gRPC сервис.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

// NewInventoryServer создаёт сервер с предзагруженными seed-данными.
func NewInventoryServer() *InventoryServer {
	now := timestamppb.Now()

	return &InventoryServer{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Прочный плазменный корпус",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart возвращает деталь по UUID.
func (s *InventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	id, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.GetUuid())
	}

	part, ok := s.parts[id]
	if !ok {
		return nil, status.Error(codes.NotFound, "деталь по заданому UUID не найдена")
	}

	return &inventoryv1.GetPartResponse{Part: part.ToContract()}, nil
}

// ListParts возвращает список деталей с опциональной фильтрацией по типу.
func (s *InventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	uuids := req.GetUuids()

	parts := make([]Part, 0)

	if len(uuids) > 0 {
		for i, id := range uuids {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s за номером %d", id, i)
			}

			part, ok := s.parts[parsedId]
			if !ok {
				return nil, status.Errorf(codes.NotFound, "деталь по UUID %s не найдена", id)
			}

			parts = append(parts, part)
		}
	} else {
		for _, part := range s.parts {
			parts = append(parts, part)
		}

		partType := req.GetPartType().String()

		switch partType {
		case inventoryv1.PartType_PART_TYPE_UNSPECIFIED.String():
		default:

			parts = slices.DeleteFunc(parts, func(part Part) bool {
				return part.PartType.String() != partType
			})
		}

		slices.SortFunc(parts, func(a, b Part) int {
			return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})
	}

	resultParts := make([]*inventoryv1.Part, len(parts))

	for i, part := range parts {
		resultParts[i] = part.ToContract()
	}

	return &inventoryv1.ListPartsResponse{Parts: resultParts}, nil
}

func (p Part) ToContract() *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      p.PartType,
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}
