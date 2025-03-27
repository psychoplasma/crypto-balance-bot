import { ExceptionFilter, Catch, ArgumentsHost, NotFoundException, HttpException, BadRequestException } from '@nestjs/common';
import { Response } from 'express';

@Catch(NotFoundException, BadRequestException, HttpException)
export class GlobalExceptionsFilter implements ExceptionFilter {
  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const status = exception.getStatus();

    response
      .status(status)
      .json({
        statusCode: status,
        message: exception.message,
      });
  }
}
