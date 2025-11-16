import { Body, Controller, Get, Param, Post } from "@nestjs/common";
import typia from "typia";

import { OrdersService } from "./orders.service";
import type { CreateOrderRequest, Order } from "../types";
import { TypiaValidationPipe } from "../pipes/typia-validation.pipe";

const createOrderAssert = typia.createAssert<CreateOrderRequest>();

const createOrderPipe = new TypiaValidationPipe<CreateOrderRequest>(
  createOrderAssert,
  "CreateOrderRequest",
);

@Controller("orders")
export class OrdersController {
  constructor(private readonly ordersService: OrdersService) {}

  @Post()
  create(
    @Body(createOrderPipe) body: CreateOrderRequest,
  ): Order {
    return this.ordersService.create(body);
  }

  @Get(":id")
  findOne(@Param("id") id: string): Order {
    return this.ordersService.findOne(id);
  }
}
