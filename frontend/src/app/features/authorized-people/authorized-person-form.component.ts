import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AuthorizedPerson, AuthorizedPersonInput, DocumentType } from '../../core/models/api.models';

@Component({
  selector: 'app-authorized-person-form',
  standalone: true,
  imports: [ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <form class="event-form person-form" [formGroup]="form" (ngSubmit)="submit()">
      <div class="form-grid">
        <label>Nombre completo
          <input formControlName="fullName" placeholder="Ej. Laura Gómez" />
          @if (invalid('fullName')) { <span class="field-error">Escribe el nombre completo.</span> }
        </label>
        <label>Tipo de documento
          <select formControlName="documentType">
            @for (type of documentTypes; track type.value) { <option [value]="type.value">{{ type.label }}</option> }
          </select>
        </label>
        <label>Número de documento
          <input formControlName="documentNumber" placeholder="Número de documento" />
          @if (invalid('documentNumber')) { <span class="field-error">Este campo es obligatorio.</span> }
        </label>
        <label>Parentesco o relación
          <input formControlName="relationship" placeholder="Ej. Madre, tío, vecino" />
          @if (invalid('relationship')) { <span class="field-error">Este campo es obligatorio.</span> }
        </label>
        <label>Teléfono de contacto
          <input formControlName="phone" placeholder="+57 300 555 0101" />
          @if (invalid('phone')) { <span class="field-error">Usa un teléfono nacional o internacional válido.</span> }
        </label>
      </div>
      <div class="form-actions">
        <button type="button" class="button button-quiet" (click)="cancel.emit()">Cancelar</button>
        <button type="submit" class="button button-primary" [disabled]="saving">{{ saving ? 'Guardando...' : (person ? 'Guardar cambios' : 'Registrar persona') }}</button>
      </div>
    </form>
  `
})
export class AuthorizedPersonFormComponent {
  private readonly formBuilder = inject(FormBuilder);
  @Input() person: AuthorizedPerson | null = null;
  @Input() saving = false;
  @Output() saved = new EventEmitter<AuthorizedPersonInput>();
  @Output() cancel = new EventEmitter<void>();
  readonly documentTypes: { value: DocumentType; label: string }[] = [
    { value: 'CC', label: 'Cédula de ciudadanía (CC)' },
    { value: 'TI', label: 'Tarjeta de identidad (TI)' },
    { value: 'CE', label: 'Cédula de extranjería (CE)' },
    { value: 'PASSPORT', label: 'Pasaporte' },
    { value: 'OTHER', label: 'Otro' }
  ];
  readonly form = this.formBuilder.nonNullable.group({
    fullName: ['', [Validators.required, Validators.maxLength(120)]],
    documentType: ['CC' as DocumentType, Validators.required],
    documentNumber: ['', [Validators.required, Validators.maxLength(40)]],
    relationship: ['', [Validators.required, Validators.maxLength(80)]],
    phone: ['', [Validators.required, Validators.pattern(/^\+?[0-9 ()-]{7,20}$/)]]
  });

  ngOnChanges(): void {
    if (this.person) {
      this.form.patchValue(this.person);
    }
  }

  invalid(control: keyof typeof this.form.controls): boolean {
    const field = this.form.controls[control];
    return field.invalid && (field.dirty || field.touched);
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saved.emit(this.form.getRawValue());
  }
}
