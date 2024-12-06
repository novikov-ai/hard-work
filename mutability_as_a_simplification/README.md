# Мутабельность как упрощение модели данных

Данный подход интересен за счет того, что можем добавлять новый функционал к уже закрытому для изменений класса. 

Многое зависит от данных. Если предполагается, что данные могут существовать независимо от класса, то, очевидно, достигается упрощение модели данных за счет того, что мы выносим нерелевантные свойства в другое место. Некоторый аналог **нормализации** баз данных. 

Однако, такой подход может и усложнить модель, так как, если эти данные должны быть неотъемлемой частью класса, то придется над ними создавать обертку, которая будет инкапсулировать под собой уже два класса.

В моем конкретном случае, например, может родиться класс PlayerWithHeight - то есть класс, обладающий свойствами обоих моделей данных.

~~~C#
using System;
using System.Collections.Generic;

public class Program
{
    public static void Main()
    {
        var player = new Player();

        // реализуем паттерн со словарем
        // ключ - хэш пользователя -> значение - новое поле
        var playersHeight = new Dictionary<int, PlayerHeight>();

        // добавление данных в словарь
        playersHeight.Add(player.GetHashCode(), new PlayerHeight(180));

        // получение данных по хэш-коду
        Console.WriteLine(playersHeight[player.GetHashCode()]);
    }
}

/// <summary>
/// Модель игрока
/// </summary>
public class PlayerHeight
{
    private int _height;
    public PlayerHeight(int height)
    {
        this._height = height;
    }

    public int Height() => _height;
}

/// <summary>
/// Модель игрока
/// </summary>
public class Player : Any
{
    private IShifter _shifter;
    private Board _board;

    public Player() {}

    public Player(IShifter shifter, Board board)
    {
        _shifter = shifter;
        _board = board;
    }

    /// <summary>
    /// Реализует логику по совершению хода игроком
    /// </summary>
    /// <precondition>Существуют возможные ходы</precondition>
    /// <postcondition>Сматченные элементы удалены, уровень заполнился новыми</postcondtion>
    public void Move(Tile[,] tiles, Position startPosition, ConsoleKey direction)
    {
        var endPosition = Coordinator.ShiftedPosition(startPosition, direction);

        var rollback = _shifter.Shift(tiles, startPosition, endPosition);
        if (!rollback.Available())
            return;

        var matchedTiles = _board.FindSets();
        if (matchedTiles.Count > 0)
        {
            foreach (var combo in matchedTiles)
            {
                foreach (var tile in combo.GetTiles())
                {
                    tiles[tile.Item2.Y, tile.Item2.X] = new Tile("✨");
                    IncrementScore(10);
                }
            }

            IncrementMoves();
        }
        else
        {
            var (beforeShiftingStart, beforeShiftingEnd) = rollback.Back();

            tiles[startPosition.Y, startPosition.X] = beforeShiftingStart;
            tiles[endPosition.Y, endPosition.X] = beforeShiftingEnd;
        }
    }

    public int Score { get; set; }
    public int Moves { get; set; }

    private void IncrementScore(int points)
    {
        Score += points;
    }

    private void IncrementMoves()
    {
        Moves++;
    }
}
~~~