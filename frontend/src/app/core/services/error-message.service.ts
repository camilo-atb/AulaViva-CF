import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ErrorMessageService {
  getMessage(error: unknown, fallback = 'No pudimos completar la solicitud.'): string {
    return error instanceof Error ? error.message : fallback;
  }
}
