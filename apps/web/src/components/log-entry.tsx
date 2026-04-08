import Link from "next/link";
import { Rating } from "@/components/rating";
import type { LogEntry as LogEntryType } from "@/lib/types";

interface LogEntryProps {
  log: LogEntryType;
  showUser?: boolean;
}

export function LogEntryItem({ log, showUser = true }: LogEntryProps) {
  return (
    <div className="border-b border-border py-2.5">
      <div className="flex items-baseline justify-between gap-2">
        <div className="min-w-0">
          <Link
            href={`/works/${log.work_slug}`}
            className="font-serif text-sm font-bold hover:text-burgundy transition-colors"
          >
            {log.work_title}
          </Link>
          {log.company_name && (
            <span className="text-xs text-muted-foreground">
              {" "}&mdash; {log.company_name}
              {log.seen_date && ` ${new Date(log.seen_date).getFullYear()}`}
            </span>
          )}
        </div>
        {log.rating && <Rating rating={log.rating} count={0} />}
      </div>
      {log.review_text && (
        <p className="text-xs text-muted-foreground mt-1 line-clamp-2 italic leading-relaxed">
          &ldquo;{log.review_text}&rdquo;
        </p>
      )}
      <div className="flex items-center gap-2 mt-1">
        {showUser && (
          <Link
            href={`/profile/${log.username}`}
            className="text-[10px] text-muted-foreground hover:text-burgundy transition-colors"
          >
            {log.username}
          </Link>
        )}
        {log.seen_date && (
          <span className="text-[10px] text-muted-foreground">
            {new Date(log.seen_date).toLocaleDateString("en-US", {
              month: "short",
              day: "numeric",
              year: "numeric",
            })}
          </span>
        )}
      </div>
    </div>
  );
}
