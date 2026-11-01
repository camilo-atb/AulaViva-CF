import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { API_BASE_URL } from '../config/api.config';
import { AuthorizedPerson, AuthorizedPersonInput, StudentDetail } from '../models/api.models';

@Injectable({ providedIn: 'root' })
export class StudentService {
  private readonly http = inject(HttpClient);

  list(search = '') {
    let params = new HttpParams();
    if (search.trim()) {
      params = params.set('search', search.trim());
    }
    return this.http.get<StudentDetail[]>(`${API_BASE_URL}/students`, { params });
  }

  get(id: number) {
    return this.http.get<StudentDetail>(`${API_BASE_URL}/students/${id}`);
  }

  listAuthorizedPeople(studentId: number) {
    return this.http.get<AuthorizedPerson[]>(`${API_BASE_URL}/students/${studentId}/authorized-people`);
  }

  createAuthorizedPerson(studentId: number, input: AuthorizedPersonInput) {
    return this.http.post<AuthorizedPerson>(`${API_BASE_URL}/students/${studentId}/authorized-people`, input);
  }

  updateAuthorizedPerson(id: number, input: AuthorizedPersonInput) {
    return this.http.put<AuthorizedPerson>(`${API_BASE_URL}/authorized-people/${id}`, input);
  }

  deleteAuthorizedPerson(id: number) {
    return this.http.delete<void>(`${API_BASE_URL}/authorized-people/${id}`);
  }
}
