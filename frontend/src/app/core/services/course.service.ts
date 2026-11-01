import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Course } from '../models/api.models';
import { API_BASE_URL } from '../config/api.config';

@Injectable({ providedIn: 'root' })
export class CourseService {
  private readonly http = inject(HttpClient);

  list() {
    return this.http.get<Course[]>(`${API_BASE_URL}/courses`);
  }
}
