import React from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string; error?: string; hint?: string;
}
export function Input({ label, error, hint, id, className = "", ...props }: InputProps) {
  const inputId = id || label?.toLowerCase().replace(/\s+/g, "-");
  return (
    <div className="w-full">
      {label && <label htmlFor={inputId} className="label">{label}{props.required && <span className="text-danger-500 ml-1">*</span>}</label>}
      <input id={inputId} className={`input ${error ? "input-error" : ""} ${className}`} aria-invalid={!!error} {...props} />
      {error && <p className="error-text" role="alert">{error}</p>}
      {hint && !error && <p className="text-xs text-gray-500 mt-1">{hint}</p>}
    </div>
  );
}

interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string; error?: string; hint?: string;
}
export function Textarea({ label, error, hint, id, className = "", ...props }: TextareaProps) {
  const inputId = id || label?.toLowerCase().replace(/\s+/g, "-");
  return (
    <div className="w-full">
      {label && <label htmlFor={inputId} className="label">{label}</label>}
      <textarea id={inputId} className={`input resize-none ${error ? "input-error" : ""} ${className}`} {...props} />
      {error && <p className="error-text">{error}</p>}
      {hint && !error && <p className="text-xs text-gray-500 mt-1">{hint}</p>}
    </div>
  );
}

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string; error?: string; options: Array<{ value: string; label: string }>;
}
export function Select({ label, error, options, id, className = "", ...props }: SelectProps) {
  const inputId = id || label?.toLowerCase().replace(/\s+/g, "-");
  return (
    <div className="w-full">
      {label && <label htmlFor={inputId} className="label">{label}</label>}
      <select id={inputId} className={`input ${error ? "input-error" : ""} ${className}`} {...props}>
        {options.map((opt) => <option key={opt.value} value={opt.value}>{opt.label}</option>)}
      </select>
      {error && <p className="error-text">{error}</p>}
    </div>
  );
}
