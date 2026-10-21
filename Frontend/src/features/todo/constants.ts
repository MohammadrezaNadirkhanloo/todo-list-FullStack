import type { Priority, TodoFilters } from "./types";

export const PRIORITY_COLORS: Record<Priority, string> = {
  low: "bg-blue-100 border-blue-200 text-blue-800",
  medium: "bg-yellow-100 border-yellow-200 text-yellow-800",
  high: "bg-red-100 border-red-200 text-red-800",
};

export const PRIORITY_LABELS: Record<Priority, string> = {
  low: "Low",
  medium: "Medium",
  high: "High",
};

export const CATEGORIES = [
  "Work",
  "Personal",
  "Shopping",
  "Health",
  "Finance",
  "Home",
  "Other",
];

export const DEFAULT_FILTERS: TodoFilters = {
  search: "",
  category: "",
  priority: "",
};