import { tags } from "typia";

export interface OrderItem {
  sku: string;
  quantity: number & tags.Type<"uint32"> & tags.Minimum<1>;
}

export interface CreateOrderRequest {
  id: string & tags.Format<"uuid">;
  customerEmail: string & tags.Format<"email">;
  items: OrderItem[];
  note?: string & tags.MaxLength<500>;
}

export interface Order {
  id: string;
  customerEmail: string;
  items: OrderItem[];
  totalItems: number;
  createdAt: string & tags.Format<"date-time">;
}
