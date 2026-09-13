import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { Moon, Sun } from "lucide-react";
import { Header } from "./components/header";
import { cn } from "./lib/utils";
import { HeroSection } from "./components/hero";
import { LogosSection } from "./components/logos-section";

export function App() {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark");
  };

  return (
    <>
      <div className="relative flex min-h-screen flex-col overflow-hidden px-4 supports-[overflow:clip]:overflow-clip">
        <Header />
        <main
          className={cn(
            "relative mx-auto max-w-4xl grow",
            // X Borders
            "before:absolute before:-inset-y-14 before:-left-px before:w-px before:bg-border",
            "after:absolute after:-inset-y-14 after:-right-px after:w-px after:bg-border",
          )}
        >
          <HeroSection />
          <LogosSection />
        </main>
      </div>
      <div className="flex min-h-svh p-6">
        <div className="flex max-w-md min-w-0 flex-col gap-4 text-sm leading-loose">
          <div>
            <h1 className="font-medium">Project ready!</h1>
            <p>You may now add components and start building.</p>
            <p>We&apos;ve already added the button component for you.</p>
            <Button className="mt-2">Button</Button>
          </div>
          <div>
            <Button variant="outline" size="icon" onClick={toggleTheme}>
              {theme === "dark" ? (
                <Sun className="h-4 w-4" />
              ) : (
                <Moon className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
