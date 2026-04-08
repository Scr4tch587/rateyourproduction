import Link from "next/link";
import { Rating } from "@/components/rating";
import type { WorkCard as WorkCardType } from "@/lib/types";

interface WorkCardProps {
  work: WorkCardType;
}

export function WorkCard({ work }: WorkCardProps) {
  return (
    <div className="border-b border-border py-2.5 flex items-start justify-between gap-4">
      <div className="min-w-0">
        <div className="flex items-baseline gap-2">
          <Link
            href={`/works/${work.slug}`}
            className="font-serif font-bold text-sm hover:text-burgundy transition-colors truncate"
          >
            {work.title}
          </Link>
          <span className="text-[10px] italic text-muted-foreground shrink-0">
            {work.type}
          </span>
        </div>
        {work.genres.length > 0 && (
          <div className="text-[11px] text-muted-foreground mt-0.5">
            {work.genres.join(" · ")}
          </div>
        )}
      </div>
      <div className="text-right shrink-0 space-y-0.5">
        <Rating rating={work.average_rating} count={work.rating_count} />
        <div className="text-[10px] text-muted-foreground">
          {work.production_count} prod.
        </div>
      </div>
    </div>
  );
}
