import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AuthorizedPersonFormComponent } from './authorized-person-form.component';

describe('AuthorizedPersonFormComponent', () => {
  let fixture: ComponentFixture<AuthorizedPersonFormComponent>;
  let component: AuthorizedPersonFormComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AuthorizedPersonFormComponent]
    }).compileComponents();
    fixture = TestBed.createComponent(AuthorizedPersonFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('does not submit incomplete data', () => {
    let submitted = false;
    component.saved.subscribe(() => submitted = true);
    component.submit();
    expect(submitted).toBeFalse();
    expect(component.form.controls.fullName.touched).toBeTrue();
  });

  it('emits a valid authorized person', () => {
    const input = {
      fullName: 'Laura Gómez',
      documentType: 'CC' as const,
      documentNumber: '123456789',
      relationship: 'Tía',
      phone: '+57 300 555 0101'
    };
    let submitted = false;
    component.saved.subscribe((value) => submitted = value.fullName === input.fullName);
    component.form.setValue(input);
    component.submit();
    expect(submitted).toBeTrue();
  });
});
