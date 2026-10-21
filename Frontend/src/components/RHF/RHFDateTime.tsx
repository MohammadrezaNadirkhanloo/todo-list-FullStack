import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/components/ui/field";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { CalendarIcon } from "lucide-react";
import { useState } from "react";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";

interface RHFDateTimeProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  required?: boolean;
  description?: string;
  placeholder?: string;
  loading?: boolean;
  disabled?: boolean;
  className?: string;
  /** فاصله‌ی دقیقه‌ها در لیست (پیش‌فرض ۵) */
  minuteStep?: number;
}

const HOURS = Array.from({ length: 12 }, (_, i) => i + 1);

/** "2026-12-31T00:00:00Z" (بدون میلی‌ثانیه) */
function toIsoZ(date: Date) {
  return date.toISOString().replace(/\.\d{3}Z$/, "Z");
}

function parseValue(value: unknown): Date | undefined {
  if (!value) return undefined;
  const d = new Date(value as string);
  return Number.isNaN(d.getTime()) ? undefined : d;
}

function RHFDateTime<T extends FieldValues>({
  name,
  label,
  required = false,
  description,
  placeholder = "MM/DD/YYYY hh:mm aa",
  loading = false,
  disabled,
  className,
  minuteStep = 5,
}: RHFDateTimeProps<T>) {
  const { control } = useFormContext<T>();
  const [isOpen, setIsOpen] = useState(false);

  const minutes = Array.from(
    { length: Math.ceil(60 / minuteStep) },
    (_, i) => i * minuteStep,
  );

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        const date = parseValue(field.value);

        // مقدار جدید را به‌صورت ISO در فرم ذخیره می‌کند
        const commit = (d: Date) => {
          d.setSeconds(0, 0);
          field.onChange(toIsoZ(d));
        };

        // اگر هنوز تاریخی انتخاب نشده، از امروز ساعت 00:00 شروع می‌کنیم
        const getBase = () => {
          if (date) return new Date(date);
          const d = new Date();
          d.setHours(0, 0, 0, 0);
          return d;
        };

        const handleDateSelect = (selected: Date | undefined) => {
          if (!selected) return;
          const next = new Date(selected);
          // ساعت و دقیقه‌ی قبلی حفظ شود
          next.setHours(date?.getHours() ?? 0, date?.getMinutes() ?? 0, 0, 0);
          commit(next);
        };

        const handleTimeChange = (
          type: "hour" | "minute" | "ampm",
          value: string,
        ) => {
          const next = getBase();
          if (type === "hour") {
            next.setHours(
              (parseInt(value) % 12) + (next.getHours() >= 12 ? 12 : 0),
            );
          } else if (type === "minute") {
            next.setMinutes(parseInt(value));
          } else {
            const isPM = next.getHours() >= 12;
            if (value === "PM" && !isPM) next.setHours(next.getHours() + 12);
            if (value === "AM" && isPM) next.setHours(next.getHours() - 12);
          }
          commit(next);
        };

        return (
          <Field
            className={cn("gap-1.5", className)}
            data-invalid={fieldState.invalid}
          >
            {label && (
              <FieldLabel htmlFor={field.name}>
                {label} :
                {required && <span className="text-destructive">*</span>}
              </FieldLabel>
            )}

            <Popover
              modal
              open={isOpen}
              onOpenChange={(open) => {
                setIsOpen(open);
                if (!open) field.onBlur();
              }}
            >
              <PopoverTrigger
                render={
                  <Button
                    id={field.name}
                    ref={field.ref}
                    type="button"
                    variant="outline"
                    disabled={disabled || loading}
                    aria-invalid={fieldState.invalid}
                    className={cn(
                      "w-full justify-start text-left font-normal",
                      !date && "text-muted-foreground",
                    )}
                  />
                }
              >
                {loading ? (
                  <Spinner className="mr-2 h-4 w-4" />
                ) : (
                  <CalendarIcon className="mr-2 h-4 w-4" />
                )}
                {date ? (
                  format(date, "MM/dd/yyyy hh:mm aa")
                ) : (
                  <span>{placeholder}</span>
                )}
              </PopoverTrigger>

              <PopoverContent className="w-auto p-0">
                <div className="sm:flex">
                  <Calendar
                    mode="single"
                    selected={date}
                    onSelect={handleDateSelect}
                    autoFocus
                  />
                  <div className="flex flex-col sm:flex-row sm:h-[300px] divide-y sm:divide-y-0 sm:divide-x">
                    {/* Hours */}
                    <ScrollArea className="w-64 sm:w-auto">
                      <div className="flex sm:flex-col p-2">
                        {HOURS.map((hour) => (
                          <Button
                            key={hour}
                            type="button"
                            size="icon"
                            variant={
                              date && date.getHours() % 12 === hour % 12
                                ? "default"
                                : "ghost"
                            }
                            className="sm:w-full shrink-0 aspect-square"
                            onClick={() =>
                              handleTimeChange("hour", hour.toString())
                            }
                          >
                            {hour}
                          </Button>
                        ))}
                      </div>
                      <ScrollBar
                        orientation="horizontal"
                        className="sm:hidden"
                      />
                    </ScrollArea>

                    {/* Minutes */}
                    <ScrollArea className="w-64 sm:w-auto">
                      <div className="flex sm:flex-col p-2">
                        {minutes.map((minute) => (
                          <Button
                            key={minute}
                            type="button"
                            size="icon"
                            variant={
                              date && date.getMinutes() === minute
                                ? "default"
                                : "ghost"
                            }
                            className="sm:w-full shrink-0 aspect-square"
                            onClick={() =>
                              handleTimeChange("minute", minute.toString())
                            }
                          >
                            {minute}
                          </Button>
                        ))}
                      </div>
                      <ScrollBar
                        orientation="horizontal"
                        className="sm:hidden"
                      />
                    </ScrollArea>

                    {/* AM / PM */}
                    <ScrollArea>
                      <div className="flex sm:flex-col p-2">
                        {["AM", "PM"].map((ampm) => (
                          <Button
                            key={ampm}
                            type="button"
                            size="icon"
                            variant={
                              date &&
                              ((ampm === "AM" && date.getHours() < 12) ||
                                (ampm === "PM" && date.getHours() >= 12))
                                ? "default"
                                : "ghost"
                            }
                            className="sm:w-full shrink-0 aspect-square"
                            onClick={() => handleTimeChange("ampm", ampm)}
                          >
                            {ampm}
                          </Button>
                        ))}
                      </div>
                    </ScrollArea>
                  </div>
                </div>
              </PopoverContent>
            </Popover>

            {description && <FieldDescription>{description}</FieldDescription>}
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        );
      }}
    />
  );
}

export default RHFDateTime;
