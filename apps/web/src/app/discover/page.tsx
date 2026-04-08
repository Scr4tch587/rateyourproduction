import { Suspense } from "react";
import { DiscoverContent } from "./discover-content";

export default function DiscoverPage() {
  return (
    <Suspense
      fallback={
        <div className="max-w-7xl mx-auto px-4 py-6">
          <h1 className="text-base font-bold mb-4">Discover</h1>
          <div className="text-xs text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <DiscoverContent />
    </Suspense>
  );
}
