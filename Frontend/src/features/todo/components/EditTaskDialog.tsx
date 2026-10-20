import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";

import { CATEGORIES } from "../constants";
import type { Priority, Task, UpdateTaskInput } from "../types";
import { buildDueDate, parseTags, splitDueDate } from "../utils/task";

type FormState = {
  title: string;
  date: string;
  time: string;
  category: string;
  priority: Priority;
  tags: string;
};

const EMPTY_FORM: FormState = {
  title: "",
  date: "",
  time: "",
  category: "",
  priority: "medium",
  tags: "",
};

const selectClass =
  "mt-1 w-full px-3 py-2 border border-input bg-background rounded-md text-sm";

type Props = {
  task: Task | null;
  isSaving: boolean;
  onClose: () => void;
  onSave: (id: number, input: UpdateTaskInput) => void;
};

export const EditTaskDialog = ({ task, isSaving, onClose, onSave }: Props) => {
  const [form, setForm] = useState<FormState>(EMPTY_FORM);

  useEffect(() => {
    if (!task) {
      setForm(EMPTY_FORM);
      return;
    }
    const { date, time } = splitDueDate(task.dueDate);
    setForm({
      title: task.title,
      date,
      time,
      category: task.category ?? "",
      priority: task.priority,
      tags: task.tags.join(", "),
    });
  }, [task]);

  const set = <K extends keyof FormState>(key: K, value: FormState[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const canSave = !!task && form.title.trim() !== "" && !isSaving;

  const handleSave = () => {
    if (!task || !canSave) return;
    onSave(task.id, {
      title: form.title.trim(),
      dueDate: buildDueDate(form.date, form.time),
      category: form.category || null,
      priority: form.priority,
      tags: parseTags(form.tags),
    });
  };

  return (
    <Dialog open={!!task} onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Task</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium">Task</label>
            <Input
              value={form.title}
              onChange={(e) => set("title", e.target.value)}
              placeholder="Task description..."
              className="mt-1"
            />
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="text-sm font-medium">Due Date</label>
              <Input
                type="date"
                value={form.date}
                onChange={(e) => set("date", e.target.value)}
                className="mt-1"
              />
            </div>
            <div>
              <label className="text-sm font-medium">Due Time</label>
              <Input
                type="time"
                value={form.time}
                disabled={!form.date}
                onChange={(e) => set("time", e.target.value)}
                className="mt-1"
              />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="text-sm font-medium">Category</label>
              <select
                value={form.category}
                onChange={(e) => set("category", e.target.value)}
                className={selectClass}
              >
                <option value="">No Category</option>
                {CATEGORIES.map((category) => (
                  <option key={category} value={category}>
                    {category}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-sm font-medium">Priority</label>
              <select
                value={form.priority}
                onChange={(e) => set("priority", e.target.value as Priority)}
                className={selectClass}
              >
                <option value="low">Low Priority</option>
                <option value="medium">Medium Priority</option>
                <option value="high">High Priority</option>
              </select>
            </div>
          </div>
          <div>
            <label className="text-sm font-medium">Tags</label>
            <Input
              value={form.tags}
              onChange={(e) => set("tags", e.target.value)}
              placeholder="Tags (comma separated)"
              className="mt-1"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={!canSave}>
            {isSaving ? "Saving..." : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};