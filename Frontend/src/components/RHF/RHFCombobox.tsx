import { useState } from "react";
import { CheckIcon, ChevronsUpDownIcon, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { cn } from "@/lib/utils";

interface Option {
  id: number;
  name: string;
  disabled?: boolean;
}

interface RHFComboboxProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  options: Option[];
  placeholder?: string;
  emptyText?: string;
  required?: boolean;
  disabled?: boolean;
  className?: string;
  multiple?: boolean;
}

function RHFCombobox<T extends FieldValues>({
  name,
  label,
  options = [],
  placeholder = "Select option...",
  emptyText = "No item found.",
  required,
  disabled,
  className,
  multiple = false,
}: RHFComboboxProps<T>) {
  const { control } = useFormContext<T>();
  const [open, setOpen] = useState(false);

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        const selectedIds: number[] = multiple
          ? Array.isArray(field.value)
            ? field.value
            : []
          : [];
        const selectedId: number | undefined = !multiple
          ? field.value
          : undefined;

        const selectedOptions = multiple
          ? options.filter((o) => selectedIds.includes(o.id))
          : options.filter((o) => o.id === selectedId);

        const isSelected = (id: number) =>
          multiple ? selectedIds.includes(id) : selectedId === id;

        const handleSelect = (item: Option) => {
          if (multiple) {
            const exists = selectedIds.includes(item.id);
            const next = exists
              ? selectedIds.filter((id) => id !== item.id)
              : [...selectedIds, item.id];
            field.onChange(next);
          } else {
            field.onChange(selectedId === item.id ? undefined : item.id);
            setOpen(false);
          }
          field.onBlur();
        };

        const removeItem = (id: number, e: React.MouseEvent) => {
          e.stopPropagation();
          field.onChange(selectedIds.filter((v) => v !== id));
        };

        return (
          <Field
            className={cn("gap-1.5", className)}
            data-invalid={fieldState.invalid || undefined}
          >
            {label && (
              <FieldLabel>
                {label}
                {required && <span className="text-destructive"> *</span>}
              </FieldLabel>
            )}

            <Popover open={open} onOpenChange={setOpen} modal>
              <PopoverTrigger asChild>
                <Button
                  type="button"
                  variant="outline"
                  disabled={disabled}
                  onClick={() => setOpen((o) => !o)}
                  className={cn(
                    "w-full justify-between font-normal h-auto min-h-9 px-3",
                    selectedOptions.length === 0 && "text-muted-foreground",
                  )}
                >
                  <div className="flex flex-wrap gap-1 py-0.5">
                    {selectedOptions.length > 0 ? (
                      multiple ? (
                        selectedOptions.map((item) => (
                          <Badge
                            key={item.id}
                            variant="secondary"
                            className="gap-1 pr-1"
                          >
                            {item.name}
                            <span
                              role="button"
                              onClick={(e) => removeItem(item.id, e)}
                              className="rounded-sm hover:bg-muted-foreground/20"
                            >
                              <XIcon size={12} />
                            </span>
                          </Badge>
                        ))
                      ) : (
                        <span>{selectedOptions[0]?.name}</span>
                      )
                    ) : (
                      <span>{placeholder}</span>
                    )}
                  </div>
                  <ChevronsUpDownIcon
                    className="opacity-50 shrink-0"
                    size={16}
                  />
                </Button>
              </PopoverTrigger>

              <PopoverContent
                align="start"
                sideOffset={4}
                className="p-0 w-[var(--radix-popover-trigger-width)] min-w-full"
              >
                <Command className="w-full">
                  <div className="sticky top-0 z-10 bg-background">
                    <CommandInput placeholder="Search..." className="h-9" />
                  </div>

                  <CommandList className="max-h-60 overflow-y-auto">
                    <CommandEmpty>{emptyText}</CommandEmpty>

                    <CommandGroup>
                      {options?.map((item) => (
                        <CommandItem
                          key={item.id}
                          value={item.name}
                          disabled={item.disabled}
                          onSelect={() => handleSelect(item)}
                        >
                          {item.name}
                          <CheckIcon
                            className={cn(
                              "ml-auto",
                              isSelected(item.id) ? "opacity-100" : "opacity-0",
                            )}
                          />
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        );
      }}
    />
  );
}

export default RHFCombobox;
