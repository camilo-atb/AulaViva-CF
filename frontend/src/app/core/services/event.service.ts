import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CreateEventInput, EventResponse } from '../models/api.models';
import { API_BASE_URL } from '../config/api.config';

@Injectable({ providedIn: 'root' })
export class EventService {
  private readonly http = inject(HttpClient);

  list() {
    return this.http.get<EventResponse[]>(`${API_BASE_URL}/events`);
  }

  get(id: number) {
    return this.http.get<EventResponse>(`${API_BASE_URL}/events/${id}`);
  }

  create(input: CreateEventInput) {
    return this.http.post<EventResponse>(`${API_BASE_URL}/events`, input);
  }
}
