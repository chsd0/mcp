package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 1. Создаем сервер
	s := server.NewMCPServer(
		"Delivery-Calculator", // Имя сервера
		"1.0.0",               // Версия
	)

	// 2. Определяем инструмент (Tool)
	// Это аналог @mcp.tool в Python
	calculateTool := mcp.NewTool(
		"calculate_delivery", // Имя функции для LLM
		mcp.WithDescription("Рассчитать стоимость доставки на основе веса и расстояния"),
		mcp.WithString("destination", mcp.Required(), mcp.Description("Город назначения")),
		mcp.WithNumber("weight_kg", mcp.Required(), mcp.Description("Вес груза в кг")),
	)

	// 3. Регистрируем обработчик (Handler)
	s.AddTool(calculateTool, calculateDeliveryHandler)

	// 4. Запускаем в режиме Stdio (стандартный ввод/вывод)
	// Агент будет запускать этот бинарник и общаться через stdin/stdout
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// Логика инструмента
func calculateDeliveryHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. Сначала приводим Arguments к карте
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		// Если аргументы не пришли или пришли не в том формате
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	// 2. Теперь безопасно достаем значения из карты
	dest, ok := args["destination"].(string)
	if !ok {
		return mcp.NewToolResultError("destination must be a string"), nil
	}

	weight, ok := args["weight_kg"].(float64)
	if !ok {
		return mcp.NewToolResultError("weight_kg must be a number"), nil
	}

	// === БИЗНЕС-ЛОГИКА ===
	basePrice := 500.0
	cost := basePrice + (weight * 50)

	resultText := fmt.Sprintf("Доставка в %s для груза %.1f кг будет стоить %.2f руб.", dest, weight, cost)

	return mcp.NewToolResultText(resultText), nil
}
