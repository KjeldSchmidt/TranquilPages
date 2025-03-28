import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { BehaviorSubject, Observable } from 'rxjs';

export interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private userSubject = new BehaviorSubject<User | null>(null);
  user$ = this.userSubject.asObservable();

  constructor(private http: HttpClient) {
    this.checkAuth();
  }

  private checkAuth() {
    this.http.get<User>(`${environment.BACKEND_URL}/api/user/me`, { withCredentials: true }).subscribe({
      next: (user) => this.userSubject.next(user),
      error: () => this.userSubject.next(null)
    });
  }

  login() {
    window.location.href = `${environment.BACKEND_URL}/auth/login`;
  }

  logout() {
    this.http.post(`${environment.BACKEND_URL}/auth/logout`, {}, { withCredentials: true }).subscribe({
      next: () => {
        this.userSubject.next(null);
        window.location.reload();
      },
      error: (error) => {
        console.error('Logout failed:', error);
        this.userSubject.next(null);
        window.location.reload();
      }
    });
  }
} 