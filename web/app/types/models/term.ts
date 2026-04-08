import type { TermDetailResponse } from "../api/term";

export interface TermDetailWithChildren extends TermDetailResponse {
  children: TermDetailWithChildren[];
}

export interface TermWithLevel extends TermDetailResponse {
  level: number;
}
