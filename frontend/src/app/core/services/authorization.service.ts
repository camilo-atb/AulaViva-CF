import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { API_BASE_URL } from '../config/api.config';
import { AuthorizationStatus, RequestResponse } from '../models/api.models';

@Injectable({ providedIn: 'root' })
export class AuthorizationService {
  private readonly http = inject(HttpClient);

  get(token: string) {
    return this.http.get<RequestResponse>(`${API_BASE_URL}/requests/${token}`);
  }

  respond(token: string, status: Exclude<AuthorizationStatus, 'pending'>) {
    return this.http.post<RequestResponse>(`${API_BASE_URL}/requests/${token}/response`, { status });
  }
}
