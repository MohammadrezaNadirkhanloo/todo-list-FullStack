import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  Calendar,
  CheckCircle,
  Circle,
  Edit,
  GripVertical,
  MoreHorizontal,
  Tag,
  Trash2,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { PRIORITY_COLORS, PRIORITY_LABELS } from "../constants";
import type { Task } from "../types";
import { formatDueDate, isOverdue } from "../utils/task";

type Props = {
  task: Task;
  sortable: boolean;
  onToggle: (task: Task) => void;
  onEdit: (task: Task) => void;
  onDelete: (id: number) => void;
};

export const SortableTaskItem = ({
  task,
  sortable,
  onToggle,
  onEdit,
  onDelete,
}: Props) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: task.id, disabled: !sortable });

  const overdue = isOverdue(task);
  const dueLabel = formatDueDate(task);

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <Card
      ref={setNodeRef}
      style={style}
      className={`transition-all duration-200 hover:shadow-md ${
        task.done ? "opacity-60" : ""
      } ${overdue ? "border-red-500 bg-red-50" : ""}`}
    >
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3 flex-1">
            {sortable && (
              <Button
                variant="ghost"
                size="icon"
                className="flex-shrink-0 cursor-grab touch-none"
                {...attributes}
                {...listeners}
              >
                <GripVertical className="h-4 w-4 text-gray-400" />
              </Button>
            )}

            <Button
              variant="ghost"
              size="icon"
              onClick={() => onToggle(task)}
              className="flex-shrink-0"
            >
              {task.done ? (
                <CheckCircle className="h-5 w-5 text-green-600" />
              ) : (
                <Circle className="h-5 w-5 text-gray-400 hover:text-gray-600" />
              )}
            </Button>

            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <p
                  className={`text-sm font-medium ${
                    task.done ? "line-through text-muted-foreground" : ""
                  }`}
                >
                  {task.title}
                </p>
                <span
                  className={`px-2 py-1 text-xs rounded-full border ${PRIORITY_COLORS[task.priority]}`}
                >
                  {PRIORITY_LABELS[task.priority]}
                </span>
                {task.category && (
                  <span className="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full">
                    {task.category}
                  </span>
                )}
              </div>

              {task.tags.length > 0 && (
                <div className="flex flex-wrap gap-1 mb-1">
                  {task.tags.map((tag, index) => (
                    <span
                      key={`${tag}-${index}`}
                      className="px-2 py-1 text-xs bg-blue-50 text-blue-700 rounded-full flex items-center gap-1"
                    >
                      <Tag className="h-3 w-3" />
                      {tag}
                    </span>
                  ))}
                </div>
              )}

              {dueLabel && (
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Calendar className="h-3 w-3" />
                  <span className={overdue ? "text-red-600 font-medium" : ""}>
                    {dueLabel}
                  </span>
                </div>
              )}
            </div>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0" />
              }
            >
              <MoreHorizontal className="h-4 w-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onEdit(task)}>
                <Edit className="mr-2 h-4 w-4" />
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onDelete(task.id)}
                className="text-red-600"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardContent>
    </Card>
  );
};
