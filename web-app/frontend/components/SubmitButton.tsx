import { useFormStatus } from 'react-dom';
import './SubmitButton.css';

export default function Submit(props: { pendingText: string, text: string }) {
  const { pending } = useFormStatus();
  return (
    <button className={'submit-btn'} disabled={pending}>
      {pending ? props.pendingText : props.text}
    </button>
  );
}
