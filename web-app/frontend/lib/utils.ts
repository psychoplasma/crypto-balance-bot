export function formatCurrency(value: bigint, decimal: bigint, currency: string): string {
  const decimalValue = value / decimal;

  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 8,
    maximumFractionDigits: 2,
  }).format(decimalValue);
}
