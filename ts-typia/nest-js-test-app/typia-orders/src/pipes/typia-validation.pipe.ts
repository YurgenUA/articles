// src/pipes/typia-validation.pipe.ts
import {
  ArgumentMetadata,
  BadRequestException,
  Injectable,
  PipeTransform,
} from "@nestjs/common";

@Injectable()
export class TypiaValidationPipe<T> implements PipeTransform {
  constructor(
    private readonly assertFn: (input: unknown) => T,
    private readonly label = "Request body",
  ) {}

  transform(value: unknown, metadata: ArgumentMetadata): T {
    if (metadata.type !== "body") {
      return value as T;
    }

    try {
      return this.assertFn(value);
    } catch (err) {
      throw new BadRequestException({
        message: `${this.label} is invalid`,
        details: String(err),
      });
    }
  }
}
