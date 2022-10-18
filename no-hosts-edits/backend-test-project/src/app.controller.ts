import { Controller, Get, Req } from '@nestjs/common';
import { Request } from 'express';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get('/some-url')
  getLocalDevMe(@Req() req: Request): Object {
    console.dir(req.rawHeaders);
    return {
      protocol: req.protocol,
      host: req.get('host'),
      url: req.url
    }
  }
  
  @Get()
  getHello(): string {
    return this.appService.getHello();
  }
}
