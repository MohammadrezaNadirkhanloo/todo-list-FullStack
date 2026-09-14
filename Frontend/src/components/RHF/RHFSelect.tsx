import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";

interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
  status?: string;
  name?: string;
  id?: number;
}

interface SelectGroup {
  groupLabel?: string;
  options: SelectOption[];
}

interface RHFSelectProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  className?: string;
  /** flat list of options — use this OR groups, not both */
  options?: SelectOption[];
  /** grouped options — use this OR options, not both */
  groups?: SelectGroup[];
}

function RHFSelect<T extends FieldValues>({
  name,
  label,
  placeholder = "Select an option",
  required = false,
  disabled = false,
  className,
  options,
  groups,
}: RHFSelectProps<T>) {
  const { control } = useFormContext<T>();

  // Normalise: if flat options given, wrap in a single unnamed group
  const resolvedGroups: SelectGroup[] =
    groups ?? (options ? [{ options }] : []);

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field
          className={cn("gap-1.5", className)}
          data-invalid={fieldState.invalid || undefined}
        >
          {label && (
            <FieldLabel htmlFor={field.name}>
              {label} :{required && <span className="text-destructive">*</span>}
            </FieldLabel>
          )}

          <Select
            value={field.value?.value ?? ""}
            onValueChange={(val) => {
              const allOptions = resolvedGroups.flatMap((g) => g.options);
              const selectedObj = allOptions.find(
                (opt) => String(opt.value) === val,
              );

              field.onChange(selectedObj || null);
              field.onBlur();
            }}
            disabled={disabled}
          >
            <SelectTrigger
              id={field.name}
              ref={field.ref}
              aria-invalid={fieldState.invalid}
            >
              <SelectValue placeholder={placeholder} />
            </SelectTrigger>

            <SelectContent position="popper">
              {resolvedGroups.map((group, gi) => (
                <SelectGroup key={gi}>
                  {group.groupLabel && (
                    <SelectLabel>{group.groupLabel}</SelectLabel>
                  )}
                  {group.options.map((opt) => (
                    <SelectItem
                      key={opt.value}
                      value={opt.value}
                      disabled={opt.disabled}
                    >
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              ))}
            </SelectContent>
          </Select>

          {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
        </Field>
      )}
    />
  );
}

export default RHFSelect;
