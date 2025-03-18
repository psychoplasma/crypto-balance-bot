'use client';

import { useFormStatus } from 'react-dom';
import { useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { signupAction } from '@/actions/actions';
import './SignUpForm.css';

export default function SignUpForm() {
  const formRef = useRef<HTMLFormElement | null>(null);
  const [error, setError] = useState('');
  const router = useRouter();

  function Submit() {
    const { pending } = useFormStatus();
    return (
      <button className={"submit-btn"} disabled={pending}>
        {pending ? 'Submitting' : 'Sign Up'}
      </button>
    );
  }

  async function submitAction(formData: FormData) {
    setError('');

    try {
      await signupAction(
        formData.get('email') as string,
        formData.get('password') as string, //FIXME: pass password hash instead
      );

      formRef.current?.reset();
      router.push('/login');
    } catch (error) {
      setError(error as string || 'Something went wrong');
    }
  };

  return (
    <form action={submitAction} ref={formRef}>
      <div className={"form-group"}>
          <label htmlFor="email">Email</label>
          <input type="email" name="email" required />
      </div>
      <div className={"form-group"}>
          <label htmlFor="password">Password</label>
          <input type="password" name="password" required />
      </div>
      {error && <p className={"error"}>{error}</p>}
      <Submit />
    </form>
  );
};
