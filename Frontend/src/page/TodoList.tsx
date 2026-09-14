import { Footer } from "@/components/footer";
import { Header } from "@/components/header";
import TodoList from "@/features/todo/TodoList";
import { cn } from "cn";

function TodoListPage() {
  return (
    <div className="relative flex min-h-screen flex-col overflow-hidden px-4 supports-[overflow:clip]:overflow-clip">
      <Header />
      <main
        className={cn(
          "relative mx-auto max-w-4xl grow h-full flex flex-col ",
          "before:absolute before:-inset-y-14 before:-left-px before:w-px before:bg-border",
          "after:absolute after:-inset-y-14 after:-right-px after:w-px after:bg-border",
        )}
      >
        <TodoList />
        <Footer />
      </main>
    </div>
  );
}

export default TodoListPage;
