package models

import "google.golang.org/protobuf/runtime/protoimpl"

type CreateProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Slug         string  `protobuf:"bytes,2,opt,name=slug,proto3" json:"slug,omitempty"`
	Image        string  `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Category     string  `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Brand        string  `protobuf:"bytes,5,opt,name=brand,proto3" json:"brand,omitempty"`
	Price        float32 `protobuf:"fixed32,6,opt,name=price,proto3" json:"price,omitempty"`
	CountInStock int32   `protobuf:"varint,7,opt,name=countInStock,proto3" json:"countInStock,omitempty"`
	Description  string  `protobuf:"bytes,8,opt,name=description,proto3" json:"description,omitempty"`
	Owner        string  `protobuf:"bytes,9,opt,name=owner,proto3" json:"owner,omitempty"`
}
type ReadProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type UpdateProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Slug         string  `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Image        string  `protobuf:"bytes,4,opt,name=image,proto3" json:"image,omitempty"`
	Category     string  `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Brand        string  `protobuf:"bytes,6,opt,name=brand,proto3" json:"brand,omitempty"`
	Price        float32 `protobuf:"fixed32,7,opt,name=price,proto3" json:"price,omitempty"`
	CountInStock int32   `protobuf:"varint,8,opt,name=countInStock,proto3" json:"countInStock,omitempty"`
	Description  string  `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	Owner        string  `protobuf:"bytes,10,opt,name=owner,proto3" json:"owner,omitempty"`
}
type DeleteProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}
type DeleteProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}
type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Slug         string  `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Image        string  `protobuf:"bytes,4,opt,name=image,proto3" json:"image,omitempty"`
	Category     string  `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Brand        string  `protobuf:"bytes,6,opt,name=brand,proto3" json:"brand,omitempty"`
	Price        float32 `protobuf:"fixed32,7,opt,name=price,proto3" json:"price,omitempty"`
	CountInStock int32   `protobuf:"varint,8,opt,name=countInStock,proto3" json:"countInStock,omitempty"`
	Description  string  `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	Rating       int32   `protobuf:"varint,10,opt,name=rating,proto3" json:"rating,omitempty"`
	NumReviews   int32   `protobuf:"varint,11,opt,name=numReviews,proto3" json:"numReviews,omitempty"`
	Owner        string  `protobuf:"bytes,12,opt,name=owner,proto3" json:"owner,omitempty"`
}
type GetAllProductsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}
type GetAllProductsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Products []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
}
