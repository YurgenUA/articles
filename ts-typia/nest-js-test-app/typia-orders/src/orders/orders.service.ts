import { Injectable, NotFoundException } from "@nestjs/common";
import { CreateOrderRequest, Order } from "../types";

@Injectable()
export class OrdersService {
  private readonly orders = new Map<string, Order>();

  create(data: CreateOrderRequest): Order {
    const totalItems = data.items.reduce(
      (sum, item) => sum + item.quantity,
      0,
    );

    const order: Order = {
      id: data.id,
      customerEmail: data.customerEmail,
      items: data.items,
      totalItems,
      createdAt: new Date().toISOString(),
    };

    this.orders.set(order.id, order);
    return order;
  }

  findOne(id: string): Order {
    const order = this.orders.get(id);
    if (!order) {
      throw new NotFoundException(`Order ${id} not found`);
    }
    return order;
  }
}
