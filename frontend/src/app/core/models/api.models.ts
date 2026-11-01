export type AuthorizationStatus = 'pending' | 'authorized' | 'not_authorized';

export interface Course {
  id: number;
  name: string;
}

export interface Student {
  id: number;
  courseId: number;
  fullName: string;
}

export interface Authorization {
  id: number;
  eventId: number;
  studentId: number;
  studentName?: string;
  requestToken: string;
  status: AuthorizationStatus;
  respondedAt?: string;
}

export interface Event {
  id: number;
  courseId: number;
  courseName?: string;
  name: string;
  description: string;
  eventDate: string;
}

export interface AuthorizationSummary {
  pending: number;
  authorized: number;
  notAuthorized: number;
}

export interface EventResponse extends Event {
  authorizations?: Authorization[];
  summary?: AuthorizationSummary;
}

export interface RequestResponse {
  authorization: Authorization;
  event: Event;
  student: Student;
}

export interface CreateEventInput {
  name: string;
  description: string;
  eventDate: string;
  courseId: number;
}

export interface ApiError {
  error?: {
    code?: string;
    message?: string;
  };
}
