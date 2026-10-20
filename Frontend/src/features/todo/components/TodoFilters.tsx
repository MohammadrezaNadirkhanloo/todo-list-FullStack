import { Search } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

import { CATEGORIES } from "../constants";
import type { Priority, TodoFilters as Filters } from "../types";

type Props = {
  value: Filters;
  onChange: (value: Filters) => void;
};

const selectClass =
  "px-3 py-2 border border-input bg-background rounded-md text-sm";

export const TodoFilters = ({ value, onChange }: Props) => (
  <Card className="mb-6">
    <CardContent className="pt-6">
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            type="text"
            value={value.search}
            onChange={(e) => onChange({ ...value, search: e.target.value })}
            placeholder="Search tasks..."
            className="pl-10"
          />
        </div>
        <select
          value={value.category}
          onChange={(e) => onChange({ ...value, category: e.target.value })}
          className={selectClass}
        >
          <option value="">All Categories</option>
          {CATEGORIES.map((category) => (
            <option key={category} value={category}>
              {category}
            </option>
          ))}
        </select>
        <select
          value={value.priority}
          onChange={(e) =>
            onChange({ ...value, priority: e.target.value as Priority | "" })
          }
          className={selectClass}
        >
          <option value="">All Priorities</option>
          <option value="low">Low Priority</option>
          <option value="medium">Medium Priority</option>
          <option value="high">High Priority</option>
        </select>
      </div>
    </CardContent>
  </Card>
);