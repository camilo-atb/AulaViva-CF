import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

export const apiErrorInterceptor: HttpInterceptorFn = (_request, next) =>
  next(_request).pipe(
    catchError((error: HttpErrorResponse) => {
      const message = error.error?.error?.message ?? 'No pudimos completar la solicitud.';
      return throwError(() => new Error(message));
    })
  );
