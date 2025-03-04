import { applyDecorators } from '@nestjs/common';
import { Transform } from 'class-transformer';
import { Timestamp } from '@google-cloud/firestore';

export function TimestampTransform() {
  return applyDecorators(
    Transform(
      ({ value }) => (value ? Timestamp.fromDate(new Date(value)) : null),
      { toClassOnly: true },
    ),
    Transform(({ value }) => (value ? value.toDate().toISOString() : null), {
      toPlainOnly: true,
    }),
  );
}
