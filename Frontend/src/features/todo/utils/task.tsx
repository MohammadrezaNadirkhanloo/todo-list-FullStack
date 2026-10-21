import { format, isValid, parseISO } from "date-fns";

import type { Task } from "../types";

export const hasTime = (value: string) => value.includes("T");

export const parseDueDate = (task: Task): Date | null => {
  if (!task.dueDate) return null;
  const date = parseISO(task.dueDate);
  return isValid(date) ? date : null;
};

export const formatDueDate = (task: Task) => {
  const date = parseDueDate(task);
  if (!date || !task.dueDate) return null;

  const dateStr = format(date, "MMM dd, yyyy");
  const timeStr = hasTime(task.dueDate) ? ` at ${format(date, "HH:mm")}` : "";
  return `${dateStr}${timeStr}`;
};

export const isOverdue = (task: Task) => {
  const date = parseDueDate(task);
  if (!date || !task.dueDate || task.done) return false;

  const deadline = hasTime(task.dueDate)
    ? date
    : new Date(date.getFullYear(), date.getMonth(), date.getDate(), 23, 59, 59);

  return deadline < new Date();
};

export const splitDueDate = (dueDate: string | null) => {
  if (!dueDate) return { date: "", time: "" };
  const parsed = parseISO(dueDate);
  if (!isValid(parsed)) return { date: "", time: "" };
  return {
    date: format(parsed, "yyyy-MM-dd"),
    time: hasTime(dueDate) ? format(parsed, "HH:mm") : "",
  };
};

export const buildDueDate = (date: string, time: string): string | null => {
  if (!date) return null;
  if (!time) return date;
  return new Date(`${date}T${time}`).toISOString();
};

export const parseTags = (value: string) =>
  value
    .split(",")
    .map((tag) => tag.trim())
    .filter((tag) => tag.length > 0);