"use client";

import { Thread } from "@/components/assistant-ui/thread";
import { ThreadList } from "@/components/assistant-ui/thread-list";
import { MyRuntimeProvider } from "./MyRuntimeProvider";

const MyApp = () => {
  return (
    <MyRuntimeProvider>
      <div className="grid h-dvh grid-cols-1 md:grid-cols-[200px_1fr] gap-x-2 px-4 py-4">
        <ThreadList />
        <Thread />
      </div>
    </MyRuntimeProvider>
  );
};
export default MyApp;
