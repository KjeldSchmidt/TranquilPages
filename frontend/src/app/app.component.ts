import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService, User } from './auth/auth.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'Welcome';
  user$: Observable<User | null>;
  
  constructor(private authService: AuthService) {
    this.user$ = this.authService.user$;
  }
  
  login() {
    this.authService.login();
  }

  logout() {
    this.authService.logout();
  }
}
