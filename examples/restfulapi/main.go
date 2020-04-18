package main

import (
	"context"
	"net/http"

	"github.com/golangpub/wine"
)

func main() {
	s := wine.NewServer()
	// Place NewBasicAuthHandler handler to do authenticating
	r := s.Router
	service := NewItemService()
	r.Get("/items/{id}", service.Get)
	r.Get("/items/list", service.List)

	r = r.Use(wine.NewBasicAuthHandler(map[string]string{"user": "password"}, "wine"))
	r.Post("/items", service.Create)

	s.Run(":8000")
}

type Item struct {
	ID    int64   `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type ItemService struct {
	items   []*Item
	counter int64
}

func NewItemService() *ItemService {
	s := new(ItemService)
	s.counter = 1
	s.items = []*Item{
		{
			ID:    1,
			Title: "Apple",
			Price: 5.8,
		},
	}
	return s
}

func (s *ItemService) Get(ctx context.Context, req *wine.Request, next wine.Invoker) wine.Responder {
	id := req.Params().Int64("id")
	for _, v := range s.items {
		if v.ID == id {
			return wine.JSON(http.StatusOK, v)
		}
	}
	return wine.Status(http.StatusNotFound)
}

func (s *ItemService) List(ctx context.Context, req *wine.Request, next wine.Invoker) wine.Responder {
	return wine.JSON(http.StatusOK, s.items)
}

func (s *ItemService) Create(ctx context.Context, req *wine.Request, next wine.Invoker) wine.Responder {
	s.counter++
	v := new(Item)
	v.ID = s.counter
	v.Title = req.Params().String("title")
	v.Price = req.Params().Float64("price")
	if v.Title == "" {
		return wine.Text(http.StatusBadRequest, "missing title")
	}
	if v.Price <= 0 {
		return wine.Text(http.StatusBadRequest, "missing price")
	}
	s.items = append(s.items, v)
	return wine.JSON(http.StatusCreated, v)
}
