/**
 * تنها منبع کلیدهای کوئریِ accessRole.
 * هیچ‌جای دیگری آرایه‌ی کلید دستی نساز.
 *
 * ساختار سلسله‌مراتبی است و تطبیق کلید prefix-based انجام می‌شود، یعنی:
 *   invalidateQueries({ queryKey: accessRoleKeys.detail(5) }) → فقط نقش ۵
 *   invalidateQueries({ queryKey: accessRoleKeys.all })       → همه‌چیز
 */
export const accessRoleKeys = {
  all: ["access-role"] as const,

  details: () => [...accessRoleKeys.all, "detail"] as const,
  detail: (id: number | string) => [...accessRoleKeys.details(), id] as const,

  /**
   * کلید جدول. این کلید در هوک جنریکِ جدول ساخته می‌شود و ما اینجا
   * داریم همان را تکرار می‌کنیم — تا وقتی در دو جا دستی نوشته شود،
   * یک غلط تایپی کافی است تا invalidate بی‌سروصدا کار نکند.
   * راه‌حل درست: هوک جدول کلیدش را export کند و اینجا import شود.
   */
  table: () => ["table-query", "/acl/roles"] as const,
} as const;
