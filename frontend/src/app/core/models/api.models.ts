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

export type DocumentType = 'CC' | 'TI' | 'CE' | 'PASSPORT' | 'OTHER';

export interface AuthorizedPerson {
  id: number;
  studentId: number;
  fullName: string;
  documentType: DocumentType;
  documentNumber: string;
  relationship: string;
  phone: string;
}

export interface StudentDetail extends Student {
  authorizedPeople: AuthorizedPerson[];
}

export interface AuthorizedPersonInput {
  fullName: string;
  documentType: DocumentType;
  documentNumber: string;
  relationship: string;
  phone: string;
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
