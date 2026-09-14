import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Mail, Phone, PhoneCallIcon } from "lucide-react";

function TelegramIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M21.5 3.5 2.7 10.8c-1.2.48-1.2 1.15-.22 1.45l4.8 1.5 1.85 5.6c.22.6.4.85.8.85.4 0 .58-.18.8-.4l1.93-1.85 4.02 2.95c.74.4 1.27.2 1.46-.68l2.64-12.4c.28-1.1-.4-1.6-1.28-1.32Zm-11.3 10.9-1.4-4.5 8.6-5.4c.42-.26.8-.12.48.16l-7.68 9.74Z" />
    </svg>
  );
}

function WhatsappIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M12.04 2C6.58 2 2.13 6.45 2.13 11.91c0 1.75.46 3.45 1.32 4.95L2 22l5.29-1.39c1.44.79 3.06 1.2 4.71 1.2h.01c5.46 0 9.9-4.45 9.9-9.91C21.91 6.45 17.5 2 12.04 2Zm5.76 14.03c-.24.68-1.38 1.3-1.9 1.36-.5.06-1 .27-3.36-.7-2.83-1.16-4.65-4-4.79-4.19-.14-.19-1.15-1.53-1.15-2.92 0-1.39.73-2.07.99-2.35.26-.28.56-.35.75-.35.19 0 .38 0 .54.01.18.01.42-.07.65.5.24.58.83 2 .9 2.15.07.15.12.32.02.51-.1.19-.15.31-.3.48-.15.17-.31.38-.44.51-.15.15-.3.31-.13.6.17.29.76 1.25 1.63 2.02 1.12.99 2.06 1.3 2.35 1.45.29.14.46.12.63-.07.17-.19.72-.84.91-1.13.19-.29.38-.24.64-.14.26.1 1.65.78 1.94.92.29.14.48.21.55.33.07.12.07.68-.17 1.36Z" />
    </svg>
  );
}

const contacts = [
  {
    label: "Mobile",
    value: "0922 203 2902",
    href: "tel:+989222032902",
    icon: Phone,
    ring: "ring-emerald-500/20",
    bg: "bg-emerald-500/10",
    text: "text-emerald-500",
  },
  {
    label: "Telegram",
    value: "@Nadirkhanloo",
    href: "https://t.me/Nadirkhanloo",
    icon: TelegramIcon,
    ring: "ring-sky-500/20",
    bg: "bg-sky-500/10",
    text: "text-sky-500",
  },
  {
    label: "WhatsApp",
    value: "@Nadirkhanloo",
    href: "https://wa.me/989222032902",
    icon: WhatsappIcon,
    ring: "ring-green-500/20",
    bg: "bg-green-500/10",
    text: "text-green-500",
  },
  {
    label: "Email",
    value: "m.nadirkhanloo1380@gmail.com",
    href: "mailto:m.nadirkhanloo1380@gmail.com",
    icon: Mail,
    ring: "ring-rose-500/20",
    bg: "bg-rose-500/10",
    text: "text-rose-500",
  },
];

export function BookACallDialog() {
  return (
    <AlertDialog>
      <AlertDialogTrigger>
        <Button variant="outline">
          <PhoneCallIcon data-icon="inline-start" /> Book a Call
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent className="max-w-sm">
        <AlertDialogHeader>
          <AlertDialogTitle className="text-center text-xl">
            Ways to Reach Me
          </AlertDialogTitle>
          <AlertDialogDescription className={"w-full"}>
            <div className="mt-2 flex flex-col gap-2 w-full">
              {contacts.map((c) => {
                const Icon = c.icon;
                return (
                  <a
                    key={c.label}
                    href={c.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className={`flex items-center gap-3 rounded-xl border p-3 ring-1 transition-colors hover:bg-muted/50 ${c.ring} w-full`}
                  >
                    <span
                      className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-full ${c.bg} ${c.text}`}
                    >
                      <Icon className="h-5 w-5" />
                    </span>
                    <span className="flex flex-col text-left">
                      <span className="text-xs text-muted-foreground">
                        {c.label}
                      </span>
                      <span className="font-medium text-foreground">
                        {c.value}
                      </span>
                    </span>
                  </a>
                );
              })}
            </div>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="w-full" variant="default">
            Got it
          </AlertDialogCancel>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
