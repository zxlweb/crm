import { z } from 'zod'

export const registerSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
  name: z.string().min(1).max(100),
  company_name: z.string().min(2).max(255),
  domain: z
    .string()
    .max(50)
    .refine((v) => v === '' || /^[a-z0-9]([a-z0-9-]{0,48}[a-z0-9])?$/.test(v), 'invalid_domain')
    .optional(),
})

export type RegisterForm = z.infer<typeof registerSchema>
