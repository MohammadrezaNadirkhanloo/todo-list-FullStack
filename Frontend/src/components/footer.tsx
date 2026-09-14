import { cn } from "@/lib/utils";
import { GithubIcon } from "@/components/icons/github-icon";
import { Logo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { FullWidthDivider } from "@/components/full-width-divider";
import { DecorIcon } from "./decor-icon";
import { Mail } from "lucide-react";

export function Footer() {
  return (
    <footer
      className={cn(
        "relative mx-auto grow flex flex-col justify-end",
        "dark:bg-[radial-gradient(35%_80%_at_15%_0%,--theme(--color-foreground/.1),transparent)]",
      )}
    >
      <div className="relative grid w-full grid-cols-6 gap-6 p-4">
        <DecorIcon className="size-4" position="top-left" />
        <DecorIcon className="size-4" position="top-right" />
        <FullWidthDivider position="top" />
        <DecorIcon className="size-4" position="bottom-left" />
        <DecorIcon className="size-4" position="bottom-right" />
        <div className="col-span-6 flex flex-col gap-4 pt-5 md:col-span-4">
          <a className="w-max" href="/">
            <Logo size={220} />
          </a>
          <p className="max-w-sm text-balance text-muted-foreground text-sm">
            A production-ready fullstack starter powered by Go, Gin, PostgreSQL,
            React, TypeScript, and shadcn/ui — everything you need to ship your
            next project faster.
          </p>
          <div className="flex gap-2">
            {socialLinks.map((item, index) => (
              <Button
                key={`social-${item.link}-${index}`}
                size="icon"
                variant="outline"
                render={<a href={item.link} target="_blank" />}
                nativeButton={false}
              >
                {item.icon}
              </Button>
            ))}
          </div>
        </div>
        <div className="col-span-3 w-full md:col-span-1">
          <span className="text-muted-foreground text-xs">Project</span>
          <div className="mt-2 flex flex-col gap-2">
            {project.map(({ href, title }) => (
              <a
                className="w-max text-sm hover:underline"
                href={href}
                key={title}
              >
                {title}
              </a>
            ))}
          </div>
        </div>
        <div className="col-span-3 w-full md:col-span-1">
          <span className="text-muted-foreground text-xs">Tech Stack</span>
          <div className="mt-2 flex flex-col gap-2">
            {techStack.map(({ href, title }) => (
              <a
                className="w-max text-sm hover:underline"
                href={href}
                key={title}
              >
                {title}
              </a>
            ))}
          </div>
        </div>
      </div>
      <FullWidthDivider />
      <div className="flex items-center justify-center gap-2 py-4">
        <p className="text-center font-light text-muted-foreground text-sm">
          &copy; {new Date().getFullYear()} Built by{" "}
          <a
            className="hover:underline"
            href="https://github.com/MohammadrezaNadirkhanloo"
            target="_blank"
            rel="noreferrer"
          >
            Mohammadreza Nadirkhanloo
          </a>
          . Open source under the MIT License.
        </p>
      </div>
    </footer>
  );
}

const project = [
  {
    title: "GitHub Repository",
    href: "https://github.com/MohammadrezaNadirkhanloo",
  },
  { title: "Documentation", href: "#" },
  {
    title: "Report an Issue",
    href: "https://github.com/MohammadrezaNadirkhanloo/issues",
  },
  { title: "License (MIT)", href: "#" },
];

const techStack = [
  { title: "Go & Gin", href: "#" },
  { title: "PostgreSQL", href: "#" },
  { title: "React & TypeScript", href: "#" },
  { title: "shadcn/ui", href: "#" },
];

const socialLinks = [
  {
    icon: <GithubIcon />,
    link: "https://github.com/MohammadrezaNadirkhanloo",
  },
  {
    icon: <Mail />,
    link: "mailto:m.nadirkhanloo1380@gmail.com",
  },
];
