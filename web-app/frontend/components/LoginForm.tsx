'use client';

import { useRouter } from 'next/navigation';
import { useRef, useState } from 'react';
import { loginAction } from '@/actions/actions';
import './LoginForm.css';
import { useFormStatus } from 'react-dom';

export default function LoginForm() {
  const formRef = useRef<HTMLFormElement | null>(null);
  const [error, setError] = useState('');
  const router = useRouter();

  function Submit() {
    const { pending } = useFormStatus();
    return (
      <button className={"submit-btn"} disabled={pending}>
        {pending ? 'Logging in...' : 'Login'}
      </button>
    );
  }

  const submitAction = async (formData: FormData) => {
    setError('');

    try {
      const { error, userId } = await loginAction(
        formData.get('email') as string,
        formData.get('password') as string, //FIXME: pass password hash instead
      );

      if (!error) {
        formRef.current?.reset();
        router.push(`/user?userId=${userId}`);
      } else {
        setError(error || 'Login failed');
      }
    } catch (error) {
      setError('An error occurred during login');
    }
  };

  return (
    <form action={submitAction}>
      <div className="form-group">
        <label htmlFor="email">Email:</label>
        <input type="email"  name="email" required/>
      </div>
      <div className="form-group">
        <label htmlFor="password">Password:</label>
        <input type="password" name="password" required/>
      </div>
      {error && <p className="error">{error}</p>}
      <Submit />
    </form>
  );
};
