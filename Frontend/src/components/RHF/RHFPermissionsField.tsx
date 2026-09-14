import { ShieldCheck } from "lucide-react";
import { useState } from "react";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";
import { useTranslation } from "react-i18next";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { formatCount } from "@/components/permissions-picker/format";
import { parsePermissionEntry } from "@/components/permissions-picker/permissions-format";
import PermissionsPicker from "@/components/permissions-picker/PermissionsPicker";
import type { PermissionsValue, Resource } from "@/components/permissions-picker/types";
import { useLocalizedResources } from "@/components/permissions-picker/use-localized-resources";
import { cn } from "@/lib/utils";

interface RHFPermissionsFieldProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  required?: boolean;
  resources: Resource[];
  className?: string;
}

function summarizeValue(value: PermissionsValue, resources: Resource[]) {
  const resourcesByKey = new Map(resources.map((r) => [r.key, r]));
  const countByResource = new Map<string, number>();

  for (const entry of value ?? []) {
    const parsed = parsePermissionEntry(entry);
    if (!parsed) continue;
    countByResource.set(parsed.resourceKey, (countByResource.get(parsed.resourceKey) ?? 0) + 1);
  }

  return {
    totalActions: value?.length ?? 0,
    resourceCount: countByResource.size,
    breakdown: Array.from(countByResource.entries()).map(([resourceKey, count]) => ({
      key: resourceKey,
      label: resourcesByKey.get(resourceKey)?.label ?? resourceKey,
      count,
    })),
  };
}

function RHFPermissionsField<T extends FieldValues>({
  name,
  label,
  required = false,
  resources,
  className,
}: RHFPermissionsFieldProps<T>) {
  const { control } = useFormContext<T>();
  const { t, i18n } = useTranslation("permissions");
  const [open, setOpen] = useState(false);
  const localizedResources = useLocalizedResources(resources);

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        const value: PermissionsValue = field.value ?? [];
        const { totalActions, resourceCount, breakdown } = summarizeValue(
          value,
          localizedResources
        );
        const hasSelection = totalActions > 0;

        return (
          <Field className="gap-1.5" data-invalid={fieldState.invalid || undefined}>
            {label && (
              <FieldLabel htmlFor={field.name}>
                {label}
                {required && <span className="text-destructive"> *</span>}
              </FieldLabel>
            )}

            <Tooltip>
              <TooltipTrigger asChild>
                <button
                  type="button"
                  id={field.name}
                  onClick={() => setOpen(true)}
                  aria-invalid={fieldState.invalid}
                  className={cn(
                    "flex h-9 w-full items-center gap-2 rounded-md border border-input bg-transparent px-3 text-sm shadow-xs outline-none transition-colors",
                    "hover:bg-accent/40 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
                    "dark:bg-input/30",
                    fieldState.invalid &&
                      "border-destructive ring-destructive/20 dark:ring-destructive/40",
                    className
                  )}
                >
                  <ShieldCheck className="size-4 shrink-0 text-muted-foreground" />
                  <span
                    className={cn(
                      "flex-1 truncate text-start",
                      !hasSelection && "text-muted-foreground"
                    )}
                  >
                    {hasSelection
                      ? t("trigger.summary", {
                          actions: formatCount(totalActions, i18n.language),
                          resources: formatCount(resourceCount, i18n.language),
                        })
                      : t("trigger.placeholder")}
                  </span>
                </button>
              </TooltipTrigger>
              {hasSelection && (
                <TooltipContent side="bottom" align="start" className="flex flex-col gap-1">
                  {breakdown.map((item) => (
                    <span key={item.key}>
                      {t("trigger.breakdownItem", {
                        resource: item.label,
                        count: formatCount(item.count, i18n.language),
                      })}
                    </span>
                  ))}
                </TooltipContent>
              )}
            </Tooltip>

            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}

            <PermissionsPicker
              open={open}
              onOpenChange={setOpen}
              resources={localizedResources}
              value={value}
              onConfirm={(next) => field.onChange(next)}
            />
          </Field>
        );
      }}
    />
  );
}

export default RHFPermissionsField;
