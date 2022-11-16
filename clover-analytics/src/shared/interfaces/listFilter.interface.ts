import { Stat } from "./stat.interface";

export interface ListFilterOptions {
  lookup: string;
  statType?: Stat;
  sportsbook?: string;
}

export interface ListSortOptions {
  sortBy?: string;
  sortOrder?: string;
  statType?: Stat;
}
