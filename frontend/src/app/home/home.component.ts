import { Component, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BookService, Book } from '../services/book.service';
import { AddBookModalComponent } from '../components/add-book-modal/add-book-modal.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, AddBookModalComponent],
  template: `
    <div class="container">
      <div class="header">
        <h1>My Books</h1>
        <button (click)="openAddBookModal()">Add Book</button>
      </div>

      <div class="books-grid">
        <div *ngFor="let book of books" class="book-card">
          <h3>{{ book.title }}</h3>
          <p class="author">by {{ book.author }}</p>
          <div class="rating">
            <span *ngFor="let i of [1,2,3,4,5]" class="star" [class.filled]="i <= book.rating">★</span>
          </div>
          <p class="comment" *ngIf="book.comment">{{ book.comment }}</p>
          <button class="delete-btn" (click)="deleteBook(book.id!)">Delete</button>
        </div>
      </div>

      <app-add-book-modal
        #addBookModal
        (bookAdded)="onBookAdded()"
      ></app-add-book-modal>
    </div>
  `,
  styles: [`
    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 2rem;
    }

    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;
    }

    button {
      padding: 0.5rem 1rem;
      background-color: #4CAF50;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 1rem;
      transition: background-color 0.2s;
    }

    button:hover {
      background-color: #45a049;
    }

    .books-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1.5rem;
    }

    .book-card {
      background: white;
      padding: 1.5rem;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    h3 {
      margin: 0 0 0.5rem 0;
      color: #333;
    }

    .author {
      color: #666;
      margin: 0 0 1rem 0;
    }

    .rating {
      margin-bottom: 1rem;
    }

    .star {
      color: #ddd;
      font-size: 1.2rem;
      margin-right: 0.2rem;
    }

    .star.filled {
      color: #ffd700;
    }

    .comment {
      color: #666;
      font-style: italic;
      margin: 0 0 1rem 0;
    }

    .delete-btn {
      background-color: #f44336;
    }

    .delete-btn:hover {
      background-color: #da190b;
    }
  `]
})
export class HomeComponent implements OnInit {
  @ViewChild('addBookModal') addBookModal!: AddBookModalComponent;
  books: Book[] = [];

  constructor(private bookService: BookService) {}

  ngOnInit() {
    this.loadBooks();
  }

  loadBooks() {
    this.bookService.getBooks().subscribe(books => {
      this.books = books;
    });
  }

  openAddBookModal() {
    this.addBookModal.open();
  }

  onBookAdded() {
    this.loadBooks();
  }

  deleteBook(id: string) {
    if (confirm('Are you sure you want to delete this book?')) {
      this.bookService.deleteBook(id).subscribe(() => {
        this.loadBooks();
      });
    }
  }
} 