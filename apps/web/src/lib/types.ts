export type WorkType = "play" | "musical" | "opera";

export interface WorkCard {
  id: string;
  slug: string;
  title: string;
  type: WorkType;
  genres: string[];
  average_rating: number;
  rating_count: number;
  weighted_score: number;
  production_count: number;
}

export interface WorkDetail {
  id: string;
  slug: string;
  title: string;
  type: WorkType;
  description?: string;
  premiere_year?: number;
  average_rating: number;
  rating_count: number;
  weighted_score: number;
  created_at: string;
  updated_at: string;
  genres: Genre[];
  creators: WorkCreator[];
  productions: Production[];
}

export interface Genre {
  id: string;
  name: string;
  slug: string;
}

export interface WorkCreator {
  person_id: string;
  name: string;
  slug: string;
  role_type: string;
}

export interface Production {
  id: string;
  work_id: string;
  slug: string;
  company_id?: string;
  company_name?: string;
  venue_id?: string;
  venue_name?: string;
  city?: string;
  country?: string;
  start_date?: string;
  end_date?: string;
  production_label?: string;
  average_rating: number;
  rating_count: number;
  weighted_score: number;
  created_at: string;
}

export interface ProductionDetail extends Production {
  work_title: string;
  work_slug: string;
  credits: ProductionCredit[];
}

export interface ProductionCredit {
  person_id: string;
  name: string;
  slug: string;
  role_type: string;
}

export interface LogEntry {
  id: string;
  user_id: string;
  work_id: string;
  production_id?: string;
  seen_date?: string;
  rating?: number;
  review_text?: string;
  created_at: string;
  work_title: string;
  work_slug: string;
  production_label?: string;
  company_name?: string;
  username: string;
}

export interface Profile {
  id: string;
  username: string;
  display_name?: string;
  avatar_url?: string;
  is_admin: boolean;
  log_count: number;
  review_count: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
}

export interface SessionProfile {
  id: string;
  username: string;
  email: string;
  display_name?: string;
  avatar_url?: string;
  is_admin: boolean;
}

export interface SignupRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface CreateLogRequest {
  work_id: string;
  production_id?: string;
  seen_date?: string;
  rating?: number;
  review_text?: string;
}

export interface SubmissionRequest {
  work_id: string;
  company_id?: string;
  venue_id?: string;
  city?: string;
  country?: string;
  start_date?: string;
  end_date?: string;
  production_label?: string;
}

export interface ProductionSubmission {
  id: string;
  work_id: string;
  work_title: string;
  submitted_by: string;
  submitted_by_name: string;
  company_id?: string;
  company_name?: string;
  venue_id?: string;
  venue_name?: string;
  city?: string;
  country?: string;
  start_date?: string;
  end_date?: string;
  production_label?: string;
  status: "pending" | "approved" | "rejected";
  notes?: string;
  created_at: string;
}
