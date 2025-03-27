'use client';

import { useRouter } from 'next/navigation';
import { useRef, useState } from 'react';
import { loginAction } from '@/actions/actions';
import { useAuth } from '@/context/AuthContext';
import SubmitButton from './SubmitButton';
import './LoginForm.css';

export default function LoginForm() {
  const { authLogin } = useAuth();
  const formRef = useRef<HTMLFormElement | null>(null);
  const [error, setError] = useState('');
  const router = useRouter();

  const submitAction = async (formData: FormData) => {
    setError('');

    try {
      const user = await loginAction(
        formData.get('email') as string,
        formData.get('password') as string,
      );
      await authLogin(user.id, user.accessToken!!);
      formRef.current?.reset();
      router.push('/user');
    } catch (error) {
      console.error('Error while login:', error);
      setError('Login failed');
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
      <SubmitButton text='Login' pendingText='Logging in...' />
    </form>
  );
};
