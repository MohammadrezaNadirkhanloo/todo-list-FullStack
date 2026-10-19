import ButtonComponent from "@/components/Button-custom";
import { FormProvider } from "@/components/RHF/FormContext";
import RHFDateTime from "@/components/RHF/RHFDateTime";
import RHFInput from "@/components/RHF/RHFInput";
import RHFSelect from "@/components/RHF/RHFSelect";
import { PlusCircle } from "lucide-react";
import { useCreateTodo } from "../api/postCreatAccessRoleApi";
import type { TodoFormData } from "./todo.schema";
import { useTodoForm } from "./useTodoForm";
import { Separator } from "@/components/ui/separator";

export function TodoForm() {
  const CATEGORIES = [
    { label: "Select Category", value: null },
    { label: "Work", value: "Work" },
    { label: "Personal", value: "Personal" },
    { label: "Shopping", value: "Shopping" },
    { label: "Health", value: "Health" },
    { label: "Finance", value: "Finance" },
    { label: "Home", value: "Home" },
    { label: "Other", value: "Other" },
  ];

  const PRIORITY_OPTIONS = [
    { label: "Medium Priority", value: null },
    { label: "Low", value: "low" },
    { label: "Medium", value: "medium" },
    { label: "High", value: "high" },
  ];

  const { CreatTodo, isLoadingSend } = useCreateTodo();
  const form = useTodoForm();

  function onSubmit(data: TodoFormData) {
    const payload = {
      title: data.title,
      dueDate: data.dueDate,
      category: data.category?.value ?? null,
      priority: data.priority?.value ?? "medium",
      tags: data.tags
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean),
    };

    CreatTodo(payload);
    form.reset();
  }

  return (
    <FormProvider form={form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="grid grid-cols-1 md:grid-cols-3 gap-3 mb-5"
      >
        <RHFInput
          name="title"
          label="New task"
          placeholder="Write a new task ..."
          required
          className="md:col-span-2"
          loading={isLoadingSend}
        />
        <ButtonComponent
          lable="Add Task"
          isLoadingSend={isLoadingSend}
          type="submit"
          className="mt-6.5"
          icon={<PlusCircle />}
        />
        <RHFDateTime name="dueDate"/>
        <RHFSelect
          name="category"
          placeholder="Select Category"
          options={CATEGORIES}
        />
        <RHFSelect
          name="priority"
          placeholder="Medium Priority"
          options={PRIORITY_OPTIONS}
        />
        <RHFInput
          name="tags"
          placeholder="Tags (comma separated)"
          required
          className="md:col-span-3"
          loading={isLoadingSend}
        />
      </form>
    </FormProvider>
  );
}
